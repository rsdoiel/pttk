import { assertGreater, assertNotEquals } from "@std/assert";
import { FmtHelp } from "./version.ts";

Deno.test(function TestFmtHelp() {
  const help_text = `%{app_name}(1) user manual | {version} {release_hash}
% Jim Doe
% {release_date}

# Name

Testing FmtHelp ...
`,
    app_name = "test_app",
    version = "v0.0.0",
    release_date = "2024-10-20",
    release_hash = 'zzyyxxww';
  let src: string;

  src = FmtHelp(help_text, app_name, version, release_date, release_hash);
  assertNotEquals(src, "");
  for (let val of [ app_name, version, release_date, release_hash ]) {
    assertGreater(src.indexOf(val), -1);
  }
});
