# Embedded browser assets

`defuddle.browser.js` is the unmodified `dist/index.full.js` bundle from
[`defuddle@0.19.2`](https://www.npmjs.com/package/defuddle/v/0.19.2), vendored
under its MIT license; see [`defuddle.LICENSE`](defuddle.LICENSE).

SHA-256: `b5a828bc49b863837931c3b59b97112a30b700b7fcd34dbe6e2a379dc1da379d`

Tap embeds the bundle so the extractor version always matches the Tap binary;
it does not write a second versioned runtime into the user's cache.
