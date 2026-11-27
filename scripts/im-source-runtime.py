"""Exercise OpenIM using two synthetic users and the recipient's history API."""

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
        raise ValueError(f"{route}: OpenIM error {result.get('errCode')}")
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
    sender, receiver = "ci-alice", "ci-bob"
    post("/user/user_register", {"users": [
        {"userID": user, "nickname": user} for user in (sender, receiver)
    ]}, admin)
    receiver_token = post("/auth/get_user_token", {"userID": receiver, "platformID": 5}, admin)["token"]
    marker = f"source-ci-{time.time_ns()}"
    post("/msg/send_msg", {
        "sendID": sender, "recvID": receiver, "senderNickname": sender,
        "sessionType": 1, "contentType": 101, "content": {"content": marker},
    }, admin)
    conversation = f"si_{sender}_{receiver}"
    for _ in range(30):
        data = post("/msg/pull_msg_by_seq", {
            "userID": receiver,
            "seqRanges": [{"conversationID": conversation, "begin": 1, "end": 16, "num": 16}],
        }, receiver_token)
        messages = data.get("msgs", {}).get(conversation, {})
        for message in messages.get("Msgs", messages.get("msgs", [])):
            content = message.get("content", "")
            if not isinstance(content, str):
                content = bytes(content).decode()
            elif marker not in content:
                import base64
                try:
                    content = base64.b64decode(content).decode()
                except (ValueError, UnicodeError):
                    continue
            if marker in content:
                (REPORTS / "private-history.json").write_text(json.dumps({
                    "assertion": "recipient history contains sent private message",
                    "conversation": conversation, "seq": message.get("seq"),
                }) + "\n")
                return
        time.sleep(1)
    raise AssertionError("recipient history never contained the sent message")


if __name__ == "__main__":
    main()
