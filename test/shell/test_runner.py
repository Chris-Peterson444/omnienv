#!/usr/bin/python3
"""Integration runner for omnienv shell behavior.

Orchestrates LXC containers via the oe binary and lxc CLI,
runs the in-container shell-test verification script, and
cleans up containers after each test.  Non-zero exit on
any failure.
"""

import json
import os
import subprocess
import shutil
import sys
import tempfile
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent.parent
OE = str(REPO_ROOT / "oe")
SHELL_TEST = REPO_ROOT / "test" / "shell" / "shell-test"

SYSTEMS = [
    "resolute",
    "noble",
    "jammy",
    "focal",
    "bionic",  # vm broken, agent not in image - quirk?
    # xenial: cloud-init status --wait forever, sys-v init handling
    # trusty: cloud-init status not supported
]


def run(cmd, **kwargs):
    print(f"run {cmd=}")
    return subprocess.run(cmd, **kwargs)


def testrun(cmd, fail_label):
    global ok
    if run(cmd).returncode != 0:
        ok = False
        print(f"FAIL: {system}/{virt} {label}", file=sys.stderr)


failed = 0
for system in SYSTEMS:
    for virt in ["container", "vm"]:
        if system == "bionic" and virt == "vm":
            # agent not present in image
            print(f"SKIP: {system}/{virt} (known broken)")
            continue

        label = f"oe-shell-test-{virt}-{system}"
        print(f"=== {system}/{virt} ===")

        with tempfile.TemporaryDirectory() as td:
            os.chdir(td)
            shutil.copy(SHELL_TEST, Path(td) / "shell-test")

            with open(".omnienv.yaml", "w") as fp:
                fp.write(f"""
                    label: oe-shell-test-{virt}
                    virtualization: {virt}
                """)

            ok = True

            testrun(
                [OE, "--launch", "-s", system, "./shell-test", "--verbose"],
               "launch/shell-test",
            )
            testrun(["lxc", "stop", label], "lxc stop")
            testrun([OE, "-s", system, "/bin/true"], "exec /bin/true")

            run(["lxc", "rm", "-f", label])

            if ok:
                print(f"PASS: {system}/{virt}")
            else:
                failed += 1

        sp = run(
            ["lxc", "image", "list", "--format", "json"],
            capture_output=True, text=True)
        if sp.returncode == 0:
            for img in json.loads(sp.stdout):
                if img.get("update_source", {}).get("alias") == system:
                    run(["lxc", "image", "delete", img["fingerprint"]])

sys.exit(failed)
