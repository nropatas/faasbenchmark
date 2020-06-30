const _ = require('lodash');

const wait = ms => new Promise((r, j) => setTimeout(r, ms));

function getDuration(startTime) {
    const end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function getSleep(event) {
    const sleep_time = _.has(event, 'extensions.request.query.sleep') ? parseInt(_.get(event, 'extensions.request.query.sleep')) : null;
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

module.exports.sleep = async (event, context) => {
    const startTime = process.hrtime();
    const params = getParameters(event);
    if (params.error) {
        return { error: params.error };
    }

    await runTest(params);

    const reused = isWarm();
    const duration = getDuration(startTime);

    return { reused, duration };
};
