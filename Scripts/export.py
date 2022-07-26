#!/usr/bin/env python3

from os import walk
import subprocess


README_HEADER = """# Discord Emoji

Discord emoji made by me (mostly for the Rust Programming Language Community Server).

## Preview

"""


generated: list[dict[str, str]] = []
for _, _, sources in walk("Sources"):
    sources.sort()

    for file in sources:
        emoji_name = file.removesuffix(".svg")
        export_name = emoji_name + ".png"

        subprocess.run(
            ["inkscape", f"--export-filename=./Export/{export_name}", "--export-width=512", f"./Sources/{file}"])
        subprocess.run(
            ["inkscape", f"--export-filename=./Preview/{export_name}", "--export-width=64", f"./Sources/{file}"])

        generated.append(
            {"name": emoji_name, "preview": f"./Preview/{export_name}", "export": f"./Export/{export_name}"})

with open("README.md", "w") as readme:
    readme.write(README_HEADER)

    for file in generated:
        readme.write("[![{name}]({preview})]({export})\n".format_map(file))
