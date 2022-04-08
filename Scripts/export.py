from dataclasses import dataclass
from pathlib import Path
from os import mkdir, walk
import subprocess


README_HEADER = """# Discord Emoji

Discord emoji made by me (mostly for the Rust Programming Language Community Server).

## Preview

"""


previews: list[dict[str, str]] = []
for _, _, files in walk("Sources"):
    files.sort()

    for file in files:
        emoji_name = file.removesuffix(".svg")
        export_name = emoji_name + ".png"

        subprocess.run(
            ["inkscape", f"--export-filename=./Export/{export_name}", "--export-width=512", f"./Sources/{file}"])
        subprocess.run(
            ["inkscape", f"--export-filename=./Preview/{export_name}", "--export-width=64", f"./Sources/{file}"])

        previews.append(
            {"name": emoji_name, "path": f"./Preview/{export_name}"})

with open("README.md", "w") as readme:
    readme.write(README_HEADER)

    for preview in previews:
        readme.write("![{name}]({path} \"{name}\")\n".format_map(preview))
