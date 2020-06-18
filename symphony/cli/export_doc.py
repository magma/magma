#!/usr/bin/env python3

import os
import re
import warnings

import pdoc
from pdoc.html_helpers import ReferenceWarning


def module_path(m: pdoc.Module, output_dir: str, ext: str):
    return os.path.join(output_dir, *re.sub(r"\.html$", ext, m.url()).split("/"))


def write_files(m: pdoc.Module, output_dir: str, **kwargs):
    f = module_path(m, output_dir, ".html")

    dirpath = os.path.dirname(f)
    if not os.access(dirpath, os.R_OK):
        os.makedirs(dirpath)

    try:
        with open(f, "w+", encoding="utf-8") as w:
            w.write("<!--\n@" + "generated\n-->\n" + m.html(**kwargs))
    except Exception:
        try:
            os.unlink(f)
        except Exception:
            pass
        raise

    for submodule in m.submodules():
        write_files(submodule, output_dir, **kwargs)


def export_doc():
    warnings.simplefilter("error", category=ReferenceWarning)
    pyinventory_module = pdoc.Module("pyinventory")
    write_files(
        pyinventory_module,
        "../docs/website/static",
        show_source_code=False,
        show_type_annotations=True,
    )


if __name__ == "__main__":
    export_doc()
