# Brief overview
- 64KB RAM
- 6502 CPU has A,X,Y register, stack pointer, PC, status register
- Instruction has different size and might require more than 1 clock cycle to complete
- The actual size of ram is 2KB but it's being viewed as 8KB by the machine, the out of range region will mirror map to the actual physical ram, so the range is being repeated by 4 times
- olc2c02 PPU
  - 8KB for pattern, stores what the graphic looks like
  - 2KB for nametable, 2D array to store the id of what needs to show in the background
  - palettes for storing color information
- Programmable ROM at the end of memory, it also contains pattern info for PPU
- Mapper will switch different bank into ROM, it stores at the same location as P-ROM, but bus will recognise it and knows it's writable
- APU, control, etc...
- Ppu memory map: https://www.nesdev.org/wiki/PPU_memory_map
  - Roughly speaking, pattern tables contain sprites (bitmap images)
  - nametable describes layout of the background
  - palette contains colour info
  - Pattern table description: https://www.nesdev.org/wiki/PPU_pattern_tables, it is an area that describe the shape of sprites in CHR rom/ram
  - Nes uses 2 bits per pixel, so 4 colors per pixel, a tile is made up of 8x8 pixel, nes has two grid side by side, each is 16x16 tiles.
  - Animated sprite works by swapping out some of the tile on the grid
  - Pixel value of 0 can be thought as transparent.
- Understand scanline and out of range scan! 
- The vertical blank ppu status register is important as it tells cpu when it's safe to update! Use for synchronising with cpu and ppu. When it reach the first out of range scaneline it'll also emit a NMI to cpu

# Running this project


# My notes
- ensure cpu is working independently, use nestest to test!! Stop till unofficial opcode
- components, implement ppu operation
- PPU is complex, do bit by bit. First make sure pattern table can be loaded and display correctly. 

# Resources
- [javidx9](https://www.youtube.com/watch?v=8XmxKPJDGU0&list=PLrOv9FMX8xJHqMvSGB_9G9nZZ_4IgteYf&index=3) series of NES emulator
- [6502 cpu datasheet](https://www.princeton.edu/~mae412/HANDOUTS/Datasheets/6502.pdf)
- [6502 guide](https://www.nesdev.org/obelisk-6502-guide/)， very complete and readable
- [6502 cpu addressing mode](https://www.nesdev.org/wiki/CPU_addressing_modes)
- [emulator test](https://www.nesdev.org/wiki/Emulator_tests) for different component
- [Nesttest](https://github.com/nwidger/nintengo/blob/master/m65go2/test-roms/nestest/nestest.log) for cpu test.
- [Status flag](https://www.nesdev.org/wiki/Status_flags) behaviour during different operation, many bugs when I'm writing my own! This is helpful reference
- [Nes architecture](https://taywee.github.io/NerdyNights/nerdynights/nesarchitecture.html) diagram and some explanation. The memory model is useful
