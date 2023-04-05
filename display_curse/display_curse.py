import json
import os
import re
import time
import xml.etree.ElementTree as ET
from argparse import ArgumentParser
from pathlib import Path


class FileWatcher:
    def __init__(self, path):
        self.ts = 0
        self.path = path

    def updated(self):
        stamp = os.stat(self.path).st_mtime
        if stamp != self.ts:
            self.ts = stamp
            return True
        return False

    def watch(self, fn):
        while True:
            if self.updated():
                print(f"{self.path} changed")
                fn()
            time.sleep(1)


def write_curses(file, out_file, gamemode):
    tree = ET.parse(file)
    root = tree.getroot()
    for node in root:
        tag = node.tag
        attrib = node.attrib

        if tag == "entry" and attrib.get("key") == gamemode:
            js = json.loads(node.text)
            curses = js.get("d", {}).get("m", [])
            curses = [re.sub(r"\[\w+\]", "", curse) for curse in curses]
            out_file.write_text("\n".join(curses))


def main():
    parser = ArgumentParser(description="Display the curses of the current run")
    parser.add_argument(
        "--input",
        "-i",
        type=Path,
        required=False,
        default=Path.home() / ".prefs" / "slice-and-dice-2",
        help="Path to the slice and  dice  save file",
    )
    parser.add_argument(
        "--output",
        "-o",
        type=Path,
        required=True,
        help="Path to output the current curse text file",
    )
    parser.add_argument(
        "--gamemode",
        "-m",
        required=False,
        default="classic",
        help="The gamemode to use",
    )
    args = parser.parse_args()
    in_file = args.input
    out_file = args.output
    gamemode = args.gamemode

    watcher = FileWatcher(in_file)
    watcher.watch(lambda:  write_curses(in_file, out_file, gamemode))

if __name__ == "__main__":
    main()
