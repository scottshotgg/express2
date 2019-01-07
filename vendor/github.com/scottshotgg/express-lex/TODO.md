# TODO

needs to have a compressor to:
- recognize scientific notation; change `[float] e [float]` to a single float
- recognize negative/positive numbers; emit unary minus and plus tokens
- recognize comma separated numbers; i.e, 12,345,678 - form of `[[int],]+`
- add `.` and allow recognization of float literals from the form `[int] . [int]`

<br>

## stages

- **lex:** <br>
  - source -> tokens<br><br>
- **parse:** <br>
  - _compress:_ <br>
    - tokens -> tokens'<br>
  - _ast:_<br>
    - tokens' -> ast<br>
  - _check:_<br>
    - type<br><br>
- transpile:<br>
    - ast -> c++<br><br>
- compiler:<br>
    - c++ -> binary<br>