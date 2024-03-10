// NOTE: this is a test function
function test() {
    console.log("test");

    const do_something = async () => {
        /** @todo make this a different function */
        await new Promise((resolve) => setTimeout(resolve, 1000));
    }

    // BUG: not awaited
    // TODO: make this awaited
    do_something();
}
