/** URL-префикс приложения, например "/game" на yacheboksarov.ru/game */
export const BASE_PATH = (process.env.REACT_APP_BASE_PATH || '').replace(/\/$/, '');

export const paths = {
  home: '/',
  login: '/login',
  play: '/play',
  results: '/results',
  debugFlags: '/debug/flags',
} as const;

export function withBasePath(path: string): string {
  if (!BASE_PATH) {
    return path;
  }
  return `${BASE_PATH}${path.startsWith('/') ? path : `/${path}`}`;
}
