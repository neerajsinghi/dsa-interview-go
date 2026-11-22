# Walkthrough generator

Every problem in the course has an animated walkthrough in `course/<module>/walkthroughs/`: a GIF that runs the solution on a real example, and a Markdown page with a mock interview and a detailed explanation. This folder holds the Python scripts that draw the GIFs.

The Go course itself still needs only the standard library. These scripts are only for regenerating the images.

## Build

```sh
pip install pillow
python3 tools/walkthroughs/build.py                    # every GIF
python3 tools/walkthroughs/build.py 04-sliding-window  # paths containing this text
```

Each run prints the size and frame count of every GIF it writes.

## How it works

- [`walkthrough.py`](walkthrough.py) is a small renderer. A `Frame` stacks widgets (`Array`, `Grid`, `LinkedList`, `Tree`, `Graph`, `KV`, `Intervals`, `Text`, `Row`) above a caption, which is what the candidate or the interviewer says during that step. `render()` lays every frame out at the same size and writes the GIF. Each frame stays on screen longer when its caption is longer.
- `tracers/<module>/<problem>.py` describes one walkthrough. It defines `TITLE`, `PATTERN` and a `frames()` generator. The tracer runs a Python port of the Go solution on the example from the problem's doc comment and yields a frame at each step, so the pictures come from a real run and are never typed in by hand.
- [`build.py`](build.py) runs every tracer and writes `course/<module>/walkthroughs/<problem>.gif`.

## Adding a walkthrough

1. Copy a tracer from a similar problem, such as [`tracers/03-two-pointers/p03_two_sum_sorted.py`](tracers/03-two-pointers/p03_two_sum_sorted.py).
2. Port the Go solution, yield a `Frame` for each meaningful step, and assert the final answer matches the doc comment.
3. Build it, then write `course/<module>/walkthroughs/<problem>.md` using the same sections as the existing pages.
4. Link it from the *Walkthrough* column of the module's problem table.

Keep captions to one or two sentences that explain *why* the step happens. Aim for 6 to 16 frames, and keep each GIF under about 300 KB.
