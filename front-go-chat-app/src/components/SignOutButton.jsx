import { useNavigate } from "react-router-dom";

export default function SignOutButton(){
  const navigate = useNavigate()
  const signOut = async() => {
    try {
      const response = await fetch('api/sign_out',{
        method: 'POST',
        credentials: 'include'
      })
      navigate('/sign_up')
    } catch (error) {
      console.log(error)
    }
    
}
    return (<>

    
  <button onClick={signOut}> Sign out </button>
  
  
  </>)
}