import config from "../config";
import withStyles from "./jss";
import { useEffect, useState } from "preact/hooks";
import { formatAbsoluteTimeWithTz } from "../util/time";

const Timer = withStyles({}, ({ classes }) => {
  const [time, setTime] = useState(Date.now());
  useEffect(() => {
    const intervalId = setInterval(() => setTime(Date.now()), 1000);
    return () => clearInterval(intervalId);
  }, []);
  if (time > config.endTime) {
    return (
      <div class="row">
        <div class={`card ${classes.card}`}>
          <div class={classes.over}>The CTF is over.</div>
        </div>
      </div>
    );
  }
  const targetEnd = time > config.startTime;
  const targetTime = targetEnd ? config.endTime : config.startTime;
  const timeLeft = targetTime - time;
  const daysLeft = Math.floor(timeLeft / (1000 * 60 * 60 * 24));
  const hoursLeft = Math.floor(timeLeft / (1000 * 60 * 60)) % 24;
  const minutesLeft = Math.floor(timeLeft / (1000 * 60)) % 60;
  const secondsLeft = Math.floor(timeLeft / 1000) % 60;
  return (
    <div class="row">
      <div class={`card ${classes.card}`}>
        <div class={classes.content}>
          <span class={classes.time}>{daysLeft}</span>
          <span class={classes.time}>{hoursLeft}</span>
          <span class={classes.time}>{minutesLeft}</span>
          <span class={classes.time}>{secondsLeft}</span>
          <span>Days</span>
          <span>Hours</span>
          <span>Minutes</span>
          <span>Seconds</span>
          <span class={classes.sub}>
            until {config.ctfName} {targetEnd ? "ends" : "starts"}
          </span>
          <span class={classes.absolute}>
            {formatAbsoluteTimeWithTz(targetTime)}
          </span>
        </div>
      </div>
    </div>
  );
});

export default Timer;
