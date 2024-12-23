import { useAuth } from "./useAuth"

 const {getUsername} = useAuth()
export const nameFormatter = (name) => (name.split('_').filter((el) => (el != getUsername())).join(' '))