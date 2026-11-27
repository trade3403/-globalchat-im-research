"""Exercise OpenIM using two synthetic users and the recipient's history API."""

import base64
import json
import os
from pathlib import Path
import time
import urllib.error
import urllib.request

BASE = "http://127.0.0.1:10002"
REPORTS = Path(os.environ["REPORTS"])


def post(route, body, token=None):
    headers = {"Content-Type": "application/json", "operationID": str(time.time_ns())}
    if token:
        headers["token"] = token
    request = urllib.request.Request(
        BASE + route, json.dumps(body).encode(), headers, method="POST"
    )
    with urllib.request.urlopen(request, timeout=5) as response:
        result = json.load(response)
    if result.get("errCode") != 0:
        raise ValueError(f"{route}: OpenIM error {result.get('errCode')} {result.get('errMsg')} {result.get('errDlt')}")
    return result.get("data") or {}


def ready_admin():
    secret = next(
        line.split(":", 1)[1].strip()
        for line in Path("config/share.yml").read_text().splitlines()
        if line.startswith("secret:")
    )
    for _ in range(60):
        try:
            return post("/auth/get_admin_token", {"secret": secret, "userID": "imAdmin"})["token"]
        except (urllib.error.URLError, ValueError, TimeoutError, KeyError):
            time.sleep(2)
    raise RuntimeError("OpenIM administrator authentication never became ready")


def main():
    admin = ready_admin()
    sender, receiver = "cialice", "cibob"
    post("/user/user_register", {"users": [
        {"userID": user, "nickname": user} for user in (sender, receiver)
    ]}, admin)
    receiver_token = post("/auth/get_user_token", {"userID": receiver, "platformID": 5}, admin)["token"]
    conversation = f"si_{sender}_{receiver}"
    sent = []
    latencies_ms = []
    for index in range(50):
        marker = f"source-ci-{time.time_ns()}-{index}"
        started = time.perf_counter()
        post("/msg/send_msg", {
            "sendID": sender, "recvID": receiver, "senderNickname": sender,
            "sessionType": 1, "contentType": 101, "content": {"content": marker},
        }, admin)
        latencies_ms.append((time.perf_counter() - started) * 1000)
        sent.append(marker)

    # A receiver token must not be allowed to pull another user's conversation.
    try:
        post("/msg/pull_msg_by_seq", {
            "userID": sender,
            "seqRanges": [{"conversationID": conversation, "begin": 1, "end": 100, "num": 100}],
        }, receiver_token)
    except (urllib.error.HTTPError, ValueError):
        unauthorized_read_rejected = True
    else:
        raise AssertionError("receiver token could read another user's message history")

    deadline = time.monotonic() + 45
    received = set()
    occurrence_count = {}
    last_seq = 0
    while time.monotonic() < deadline and len(received) < len(sent):
        data = post("/msg/pull_msg_by_seq", {
            "userID": receiver,
            "seqRanges": [{"conversationID": conversation, "begin": 1, "end": 256, "num": 256}],
        }, receiver_token)
        messages = data.get("msgs", {}).get(conversation, {})
        for message in messages.get("Msgs", messages.get("msgs", [])):
            content = message.get("content", "")
            if not isinstance(content, str):
                content = bytes(content).decode()
            try:
                decoded = json.loads(content)
            except (ValueError, TypeError):
                try:
                    decoded = json.loads(base64.b64decode(content).decode())
                except (ValueError, UnicodeError):
                    continue
            text = decoded.get("content", "")
            if text in sent:
                occurrence_count[text] = occurrence_count.get(text, 0) + 1
                received.add(text)
                last_seq = max(last_seq, int(message.get("seq", 0)))
        if len(received) < len(sent):
            time.sleep(1)

    missing = sorted(set(sent) - received)
    if missing:
        raise AssertionError(f"recipient history missed {len(missing)} of {len(sent)} messages")
    latencies_ms.sort()
    percentile = lambda q: latencies_ms[min(len(latencies_ms)-1, int((len(latencies_ms)-1)*q))]
    result = {
        "assertion": "50 private messages persisted and were readable by recipient",
        "sent": len(sent),
        "received_unique": len(received),
        "missing": len(missing),
        "duplicate_count_in_history": sum(count - 1 for count in occurrence_count.values() if count > 1),
        "unauthorized_history_read_rejected": unauthorized_read_rejected,
        "send_rest_latency_ms": {
            "p50": round(percentile(.50), 3),
            "p95": round(percentile(.95), 3),
            "p99": round(percentile(.99), 3),
        },
        "last_sequence": last_seq,
        "note": "REST request latency is not WebSocket end-to-end delivery latency.",
    }
    (REPORTS / "private-history.json").write_text(json.dumps(result, indent=2) + "\n")



if __name__ == "__main__":
    main()
