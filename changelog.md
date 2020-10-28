## Changelog

### 1.1.1 - 28/8/2020 - Jose Attento (jose.attento@gmail.com)
- Add a default of 64 value for bitmaps length, the value is assumed if no length is indicated.

### 1.1.0 - 28/8/2020 - Jose Attento (jose.attento@gmail.com)
- Replace suport to go native types (as fields) with the new field types binary, llbinary and lllbinary.
- Add new input validation to bitmap type.
- Add MasterCard ISO87 struct template.
- Add more documentation, now all exported elements are documented.

### 1.0.1 - 26/8/2020 - Jose Attento (jose.attento@gmail.com)
- Refactor code inencode.go and decode.go with suggestions from https://codeclimate.com/github/jattento/go-iso8583
- Add new test cases
- Improve error messages

### 1.0.0 - 16/8/2020 - Jose Attento (jose.attento@gmail.com)
- Add ISO-8583 Marshal and Unmarshal.
