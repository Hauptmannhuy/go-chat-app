import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import SignInInput from "./SignInInput"
import SignUpInput from "./SignUpInput"


export default function AuthenticationForm(){
  const type = window.location.pathname
  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [pass, setPass] = useState("")

  const [authSuccess, setAuthSuccess] = useState(false)
  const [authStatus, setAuthStatus] = useState(false)
  const [errorStatus, setErrorStatus] = useState(false)
  const [errorBody, setErrorBody] = useState(null)

  const navigate = useNavigate()

  useEffect(() => {
    async function auth() {
      if (authStatus) {
        const message = type == '/sign_up' ? {"username": name, "email": email, "password": pass} : {"username": name, "password": pass}
        try {
          const response = await fetch(`${'api'+type}`, {
            method: "POST",
            mode: "cors",
            body: JSON.stringify(message),
            credentials: 'include'
          });

          if (response.status == 200) {
            setAuthSuccess(true);
          } else if (response.status == 401) {
            setErrorStatus(true)
            const json = await response.json()
            setErrorBody(json)
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
  
  useEffect(() => {
    if (errorStatus && errorBody) {
      console.log(errorBody);
    }
  }, [errorStatus, errorBody]);
    
  return (
    <>
    
    {errorStatus ? (<div>{errorBody}</div>) : null} 
    

    {type == '/sign_up' ? (
    <SignUpInput
      designation={type}
      nameHandler={setName}
      emailHandler={setEmail}
      passHandler={setPass}
      actionSubmit={setAuthStatus}/>
    ) : (
    <SignInInput
      designation={type}
      nameHandler={setName}
      passHandler={setPass}
      actionSubmit={setAuthStatus}
      />
      )}
  
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

