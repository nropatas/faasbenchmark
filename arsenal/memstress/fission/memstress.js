const url = require('url');
const _ = require('lodash');

const MEGABYTE = 1024 * 1024;

function memIntensive(level, memSize) {
    const available_memory = parseInt(memSize, 10);
    const amountInMB = available_memory - (available_memory / 10) * (4 - level);
    console.log(amountInMB);
    Buffer.alloc(amountInMB * MEGABYTE, 'a');
}

function getQueryParams(context) {
    const urlParts = url.parse(context.request.url, true);
    return urlParts.query;
}

function getDuration(startTime) {
    const end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function getLevel(event) {
    const intensityLevel = event.level ? parseInt(event.level) : null;
    if (!intensityLevel || intensityLevel < 1) {
        return {"error": "invalid level parameter"};
    }
    return intensityLevel;
}

function getParameters(event) {
    return getLevel(event);
}

function runTest(intensityLevel, memSize){
    memIntensive(intensityLevel, memSize);
}

function isWarm() {
    const is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

function getMemSize(context) {
    const funcName = _.find(context.request.headers, (value, key) => key.toLowerCase() === 'x-fission-function-name');
    return _.chain(funcName)
        .split('-')
        .last()
        .replace('mb', '')
        .value();
}

module.exports = async (context) => {
    const query = getQueryParams(context);
    const memSize = getMemSize(context);

    const startTime = process.hrtime();
    const params = getParameters(query);
    if (params.error) {
        return {
            status: 200,
            body: { error: params.error },
        };
    }

    runTest(params, memSize);

    const reused = isWarm();
    const duration = getDuration(startTime);

    return {
        status: 200,
        body: { reused, duration },
    };
};
