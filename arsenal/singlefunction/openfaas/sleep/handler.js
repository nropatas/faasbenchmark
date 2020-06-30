const wait = ms => new Promise((r, j) => setTimeout(r, ms));

function getDuration(startTime) {
    const end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function getSleep(event) {
    const sleep_time = event.query.sleep ? parseInt(event.query.sleep) : null;
    if (!sleep_time && sleep_time !== 0) {
        return {"error": "invalid sleep parameter"};
    }
    return sleep_time;
}

function getParameters(event) {
    return getSleep(event);
}

async function runTest(sleep_time){
    await wait(sleep_time);
}

function isWarm() {
    const is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

module.exports = async (event, context) => {
    const startTime = process.hrtime();
    const params = getParameters(event);
    if (params.error) {
        return context.status(200).succeed({
            error: params.error,
        });
    }

    await runTest(params);

    const reused = isWarm();
    const duration = getDuration(startTime);

    return context.status(200).succeed({
        reused: reused,
        duration: duration,
    });
}
