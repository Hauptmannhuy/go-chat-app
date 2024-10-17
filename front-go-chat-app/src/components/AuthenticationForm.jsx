import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import AuthInput from "./AuthInput"


export default function AuthenticationForm(){
  const type = window.location.pathname
  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [pass, setPass] = useState("")

  const [authSuccess, setAuthSuccess] = useState(false)
  const [authStatus, setAuthStatus] = useState(false)

  const navigate = useNavigate()

  useEffect(() => {
    async function auth() {
      if (authStatus) {
        const input = [name,email,pass]
        const message = {"username": input[0], "email": input[1], "password": input[2]}
        try {
          const response = await fetch(`${'api'+type}`, {
            method: "POST",
            mode: "cors",
            body: JSON.stringify(message),
            credentials: 'include'
          });

          if (response.status == 200) {
            setAuthSuccess(true);
          }

        } catch (error){
          console.log(`${type} failed:`, error);
        } finally {
          setAuthStatus(false);
        }
      }
    }
    
    auth()
  }, [authStatus, name, email, pass])


  useEffect(() => {
    if (authSuccess) navigate('/chats');
  }, [authSuccess]);
  
  
  console.log(type)
  
  return (
    <>

    <AuthInput
    designation={type}
    nameHandler={setName}
    emailHandler={setEmail}
    passHandler={setPass}
    actionSubmit={setAuthStatus}/>
  
    <input type="checkbox"/>
    Remember me
  

    {
      type != '/sign_up' ? (
      <div>
        <p>Don't have an account?</p>
        <button onClick={() => navigate('/sign_up')}>Sign Up </button>
      </div>
      ) 
      : (
        <div>
        <p>Have an account?</p>
        <button onClick={() => navigate('/sign_in')}>Sign In </button>
      </div>
        )
    }
  
  
  </>
  )
}

