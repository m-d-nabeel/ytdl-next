"use client";
import { useState, useEffect } from "react";

const Timer = () => {
  const [time, setTime] = useState(900); // 900 seconds = 15 minutes

  useEffect(() => {
    const interval = setInterval(() => {
      setTime((prevTime) => prevTime - 1);
    }, 1000);

    return () => {
      clearInterval(interval);
    };
  }, []);

  useEffect(() => {
    if (time === 0) {
      setTime(900); // Reset the timer to 15 minutes
    }
  }, [time]);

  const formatTime = (seconds: number) => {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes.toString().padStart(2, "0")}:${remainingSeconds
      .toString()
      .padStart(2, "0")}`;
  };

  return (
    <div>
      <h1>Timer: {formatTime(time)}</h1>
    </div>
  );
};

export default Timer;
