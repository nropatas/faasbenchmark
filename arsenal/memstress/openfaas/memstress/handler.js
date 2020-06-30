const MEGABYTE = 1024 * 1024;

function memIntensive(level) {
    const available_memory = parseInt(process.env.function_memory_size, 10);
    const amountInMB = available_memory - (available_memory / 10) * (4 - level);
    console.log(amountInMB);
    Buffer.alloc(amountInMB * MEGABYTE, 'a');
}

function getDuration(startTime) {
    const end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function getLevel(event) {
    const intensityLevel = event.query.level ? parseInt(event.query.level) : null;
    if (!intensityLevel || intensityLevel < 1) {
        return {"error": "invalid level parameter"};
    }
    return intensityLevel;
}

function getParameters(event) {
    return getLevel(event);
}

function runTest(intensityLevel){
    memIntensive(intensityLevel);
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

    runTest(params);

    const reused = isWarm();
    const duration = getDuration(startTime);

    return context.status(200).succeed({
        reused: reused,
        duration: duration,
    });
}
