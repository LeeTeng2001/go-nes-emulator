# Brief overview
- 64KB RAM
- 6502 CPU has A,X,Y register, stack pointer, PC, status register
- Instruction has different size and might require more than 1 clock cycle to complete
- The actual size of ram is 2KB but it's being viewed as 8KB by the machine, the out of range region will mirror map to the actual physical ram, so the range is being repeated by 4 times
- PPU
  - 8KB for pattern, stores what the graphic looks like
  - 2KB for nametable, 2D array to store the id of what needs to show in the background
  - palettes for storing color information
- Programmable ROM at the end of memory, it also contains pattern info for PPU
- Mapper will switch different bank into ROM, it stores at the same location as P-ROM, but bus will recognise it and knows it's writable
- APU, control, etc...

# Resources
- [javidx9](https://www.youtube.com/watch?v=8XmxKPJDGU0&list=PLrOv9FMX8xJHqMvSGB_9G9nZZ_4IgteYf&index=3) series of NES emulator
- [6502 cpu datasheet](https://www.princeton.edu/~mae412/HANDOUTS/Datasheets/6502.pdf)
- [6502 guide](https://www.nesdev.org/obelisk-6502-guide/)， very complete and readable
- [6502 cpu addressing mode](https://www.nesdev.org/wiki/CPU_addressing_modes)
- [emulator test](https://www.nesdev.org/wiki/Emulator_tests) for different component
- [Nesttest](https://github.com/nwidger/nintengo/blob/master/m65go2/test-roms/nestest/nestest.log) for cpu test.
- [Status flag](https://www.nesdev.org/wiki/Status_flags) behaviour during different operation, many bugs when I'm writing my own! This is helpful reference
- [Nes architecture](https://taywee.github.io/NerdyNights/nerdynights/nesarchitecture.html) diagram and some explanation. The memory model is useful
