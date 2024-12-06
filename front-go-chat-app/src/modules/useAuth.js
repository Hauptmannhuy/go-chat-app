export const useAuth = () => {
  const userAuthenticated = () => {

    if (document.cookie != '') return false;
    return true;
  }

  const getUsername = () =>  document.cookie.split('=')[1]

  
  return { userAuthenticated, getUsername}
}