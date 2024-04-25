"use client";

export function getCookieValue(value: string) {
  return document.cookie
    .split(";s")
    .find((row) => row.startsWith(`${value}=`))
    ?.split("=")[1];
}

export function cookies() {
  return document.cookie;
}
