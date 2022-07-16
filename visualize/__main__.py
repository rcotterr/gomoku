import pygame


def main():
    # imports the pygame library module
    # initilize the pygame module
    pygame.init()
    # Setting your screen size with a tuple of the screen width and screen height
    display_screen = pygame.display.set_mode((1600, 1040))
    # Setting a random caption title for your pygame graphical window.
    pygame.display.set_caption("Gomoku")
    # Update your screen when required
    pygame.display.update()
    display_screen.fill((222, 184, 135))

    width = 54
    lines = 19
    start = 40
    plus = 32
    x_start, y_start = start, start
    x_end, y_end = x_start + width * (lines - 1), y_start + width * (lines - 1)
    lines = 19

    for i in range(lines):  # for x
        pygame.draw.lines(display_screen, (0, 0, 0), True, ((x_start, y_start), (x_start, y_end)))
        x_start += width

    x_start = start
    for i in range(lines):  # for y
        pygame.draw.lines(display_screen, (0, 0, 0), True, ((x_start, y_start), (x_end, y_start)))
        y_start += width

    pygame.display.update()

    black_stone_image = pygame.image.load("radio-button.png")
    clock = pygame.time.Clock()

    # quit the pygame initialization and module
    running = True

    while running:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                running = False
            elif event.type == pygame.MOUSEBUTTONUP or event.type == pygame.MOUSEBUTTONDOWN:
                print(event)
                x, y = event.pos
                display_screen.blit(black_stone_image, (x-plus, y-plus))
                pygame.display.update()
                # print(event.x, event.y)
                # print(event.flipped)
                # print(event.which)
                # can access properties with
                # proper notation(ex: event.y)
            clock.tick(60)

        display_screen.blit(black_stone_image, (300, 500))
        pygame.display.update()


    pygame.quit()
    # End the program
    quit()
    pass


if __name__ == "__main__":
    main()
