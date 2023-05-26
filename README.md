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
- Nametable (https://www.nesdev.org/wiki/PPU_nametables), quite complicated, most complex part of emulator:
  - Single nametable is (32x8) x (32x8), we have 4 nametable, one pattern table can fit 1/4 of single nametable. 
  - The last 2 row of nametable is used as attribute table
  - That's how nintendo get its display resolution! 32x8=256 x (32-2)x8=240
  - Actual VRAM only contains 2 nametable but there's 4 addressable nametable.
  - Same as before, some register controls mirroring of X or Y of nametable. Look at the link for visualisation
  - Scroll register controls scrolling to another nametable.
  - Offset can be cleverly calculated by combining bits! Nintendo is genius
  - The last 2 rows is 32 * 2 bytes = 64, so we can split up the nametable into 8x8 region and each attribute defines the color palette to use. Each region can be further divde into 2x2 sub-region, each sub-region pallette is controlled by 2 bits in single atrribute byte 
  - Render cycle diagram: https://www.nesdev.org/w/images/default/4/4f/Ppu.svg, when rendering current tile prepare for next one during different cycle
- Loopy rendering (widely used in NES emulator), he's the guy that define a convinient way for us to setup this rendering pipeline in emulator
- Input control, it's relatively simple, NES only has 8 buttons (which can map to single byte): https://www.nesdev.org/wiki/Input_devices
  - It's a parallel write, serial read operation. The controller is memory mapped at cpu bus.
  - Write to $4016 or %4017 to capture snapshot of controller (1/2) state
  - Read polled data one bit at a time from $4016 or $4017 starting from MSB
  - Supports two controller
- Sprites: https://www.nesdev.org/wiki/PPU_OAM
  - Where is sprite stored? In object attribute memory (special ram in ppu, not accessible directly)
  - OAM is 256 bytes, each sprite has (x, y, tileID, attribute info) that takes up 4 bytes,
  - A total of 64 sprites can be stored
  - Support two types of sprite, 8x8 and 8x16
  - Being accessed by CPU via OAM reg (one write because OAM address range is 256) to PPU
  - If we need to populate whole sprite area via OAM reg needs 256 rw, too slow
  - Mystical register $4014 can be write to, at this point CPU is suspended, DMA kicks in to write whole pages from CPU 256 times to OAM
  - 4 times faster than writing OAM reg
  - Render process: When scanline reach the end, it evaluates which sprite to render in next frame (by comparing y axis)
  - rendering along a scanline increase cycle while decrease sprite X, so if sprite x == 0 we should start to render the sprite
  - Nes can only render maximum of 8 sprites in a given scanline, if we got more than that then a sprite overflow flag should be set
  - The lowest bit sprite has the maximum priority if overlap.

# Running this project


# My notes
- ensure cpu is working independently, use nestest to test, the old one contains some bug remember to find the updated file!! Stop till unofficial opcode
- components, implement ppu operation
- PPU is complex, do bit by bit. First make sure pattern table can be loaded and display correctly. Only need some bits in status register and control register in order to get it working

# Resources
- [javidx9](https://www.youtube.com/watch?v=8XmxKPJDGU0&list=PLrOv9FMX8xJHqMvSGB_9G9nZZ_4IgteYf&index=3) series of NES emulator
- [6502 cpu datasheet](https://www.princeton.edu/~mae412/HANDOUTS/Datasheets/6502.pdf)
- [6502 guide](https://www.nesdev.org/obelisk-6502-guide/)， very complete and readable
- [6502 cpu addressing mode](https://www.nesdev.org/wiki/CPU_addressing_modes)
- [emulator test](https://www.nesdev.org/wiki/Emulator_tests) for different component
- [Nesttest](https://github.com/nwidger/nintengo/blob/master/m65go2/test-roms/nestest/nestest.log) for cpu test.
- [Status flag](https://www.nesdev.org/wiki/Status_flags) behaviour during different operation, many bugs when I'm writing my own! This is helpful reference
- [Nes architecture](https://taywee.github.io/NerdyNights/nerdynights/nesarchitecture.html) diagram and some explanation. The memory model is useful
