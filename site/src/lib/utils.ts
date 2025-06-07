export function debounce(callback: () => void, delay = 200): () => void {
  let timeout: ReturnType<typeof setTimeout> | null = null;

  return () => {
    if (timeout) clearTimeout(timeout);
    timeout = setTimeout(callback, delay);
  };
}
