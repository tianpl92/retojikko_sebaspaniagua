# Schema Analysis Report

Note: `runs.jsonl` labels every row as `spec: "a"`, so I inferred the actual spec tier from the output filename suffix (`_a`, `_b`, `_c`) to build the requested comparison matrix.

## Individual schema scores

| schema file | spec name | model | score |
|---|---:|---|---:|
| r1_gpt5.5_a.sql | spec A | gpt5.5 | 48/100 |
| r1_nemotron_a.sql | spec A | nemotron | 46/100 |
| r1_owlalpha_a.sql | spec A | owlalpha | 54/100 |
| r2_gpt5.5_b.sql | spec B | gpt5.5 | 62/100 |
| r2_nemotron_b.sql | spec B | nemotron | 66/100 |
| r2_owlalpha_b.sql | spec B | owlalpha | 58/100 |
| r3_gpt5.5_c.sql | spec C | gpt5.5 | 83/100 |
| r3_nemotron_c.sql | spec C | nemotron | 76/100 |
| r3_owlalpha_c.sql | spec C | owlalpha | 80/100 |

## 3x3 matrix: spec vs model

| type spec | gpt5.5 | nemotron | owlalpha | Spread (avg pts) |
|---|---:|---:|---:|---:|
| Spec A (minimal) | 48 | 46 | 54 | 49.3 |
| Spec B (balanced) | 62 | 66 | 58 | 62.0 |
| Spec C (comprehensive) | 83 | 76 | 80 | 79.7 |

## Brief evaluation notes

- Spec A schemas are structurally sound but overbuilt for a minimal template: they ship 6 tables, no soft delete, no comments, and no indexes.
- Spec B schemas add soft delete and better integrity/index coverage, but still drift beyond a balanced template and remain undocumented.
- Spec C schemas are the strongest overall, with comments and broader indexing, but two files have self-contained SQL issues that reduce the score:
  - `r2_owlalpha_b.sql` uses `gen_random_uuid()` without creating the required extension.
  - `r3_nemotron_c.sql` creates `citex` instead of `citext`.
  - `r3_owlalpha_c.sql` uses `CITEXT` without enabling the extension.

## Overall takeaway

The comprehensive spec (C) produced the best schemas by a clear margin. The minimal and balanced specs were followed less strictly, with all models tending to generate fuller e-commerce schemas than the spec tier suggests. The main quality differentiators were documentation coverage, self-contained SQL correctness, and index strategy.
