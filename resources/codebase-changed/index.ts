// NOTE: this function is just for testing
function test() {
    console.log("test");

    /** @todo make this a different function */
    const do_something = async () => {
        await new Promise((resolve) => setTimeout(resolve, 1000));
    }

    // BUG: not awaited
    // TODO: make this awaited
    do_something();
}
