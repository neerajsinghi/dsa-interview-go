"""Build the walkthrough GIFs.

Every tracer in tracers/<module>/<problem>.py defines TITLE, PATTERN and a
frames() function. Its GIF is written to course/<module>/walkthroughs/<problem>.gif.

    python3 tools/walkthroughs/build.py                  # build everything
    python3 tools/walkthroughs/build.py 04-sliding       # only paths containing this text
"""

import importlib.util
import pathlib
import sys
import traceback

HERE = pathlib.Path(__file__).resolve().parent
ROOT = HERE.parent.parent
sys.path.insert(0, str(HERE))

import walkthrough  # noqa: E402


def load(path):
    spec = importlib.util.spec_from_file_location(f"tracer_{path.parent.name}_{path.stem}", path)
    mod = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(mod)
    return mod


def main(filters):
    tracers = sorted((HERE / "tracers").glob("*/p*.py"))
    if filters:
        tracers = [t for t in tracers if any(f in str(t) for f in filters)]
    failed = 0
    for t in tracers:
        out = ROOT / "course" / t.parent.name / "walkthroughs" / f"{t.stem}.gif"
        try:
            mod = load(t)
            frames = list(mod.frames())
            walkthrough.render(mod.TITLE, mod.PATTERN, frames, str(out))
            print(f"{out.stat().st_size // 1024:5d} KB  {len(frames):3d} frames  {out.relative_to(ROOT)}")
        except Exception:
            failed += 1
            print(f"FAILED {t.relative_to(ROOT)}")
            traceback.print_exc()
    print(f"{len(tracers) - failed} built, {failed} failed")
    return 1 if failed else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
