import json
import re
import time
import xml.etree.ElementTree as ET
from argparse import ArgumentParser
from pathlib import Path


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

    with open(in_file, "r") as in_fd:
        while True:
            tree = ET.parse(in_fd)
            root = tree.getroot()
            for node in root:
                tag = node.tag
                attrib = node.attrib

                if tag == "entry" and attrib.get("key") == gamemode:
                    js = json.loads(node.text)
                    curses = js.get("d", {}).get("m", [])
                    curses = [re.sub(r"\[\w+\]", "", curse) for curse in curses]
                    out_file.write_text("\n".join(curses))
                    break

            in_fd.seek(0)
            time.sleep(1)


if __name__ == "__main__":
    main()
