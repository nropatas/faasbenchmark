const wait = ms => new Promise((r, j) => setTimeout(r, ms));

function getDuration(startTime) {
    const end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function getSleep(event) {
    const sleep_time = event.sleep ? parseInt(event.sleep) : null;
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

exports.sleep = async (args) => {
    const startTime = process.hrtime();
    const params = getParameters(args);
    if (params.error) {
        return {
            statusCode: 200,
            headers: { 'Content-Type': 'application/json' },
            body: { error: params.error },
        };
    }

    await runTest(params);

    const reused = isWarm();
    const duration = getDuration(startTime);

    return {
        statusCode: 200,
        headers: {
            'Content-Type': 'application/json',
        },
        body: { reused, duration },
    }
}
