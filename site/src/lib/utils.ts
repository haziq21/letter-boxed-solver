export function debounce(callback: () => void, delay = 200): () => void {
  let timeout: ReturnType<typeof setTimeout> | null = null;

  return () => {
    if (timeout) clearTimeout(timeout);
    timeout = setTimeout(callback, delay);
  };
}

export function arrEq<T>(a: T[], b: T[]): boolean {
  if (a.length !== b.length) return false;
  return a.every((val, i) => val === b[i]);
}
