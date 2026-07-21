export const setToken = (token: string, remember: boolean) => {
  if (remember) {
    localStorage.setItem('token', token);
    sessionStorage.removeItem('token');
  } else {
    sessionStorage.setItem('token', token);
    localStorage.removeItem('token');
  }
};

export const getToken = () => {
  return localStorage.getItem('token') || sessionStorage.getItem('token');
};

export const clearToken = () => {
  localStorage.removeItem('token');
  sessionStorage.removeItem('token');
};
