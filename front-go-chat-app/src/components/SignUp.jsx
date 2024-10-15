import { useEffect, useState } from "react"
import { Link } from "react-router-dom";
import { useNavigate } from "react-router-dom";


function signUpForm(){
  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [pass, setPass] = useState("")

  const navigate = useNavigate()


  const [data,setData] = useState(null)
  const [signUpStatus, invokeSignUp] = useState(false)

  console.log(data)

  useEffect(() => {
    async function signUp() {
      if (signUpStatus) {
        const input = [name,email,pass]
        const message = {"username": input[0], "email": input[1], "password": input[2]}
        try {
        const response = await fetch(`api/sign_up`, {
            method: "POST",
            mode: "cors",
            body: JSON.stringify(message),
            credentials: 'include'
          })
          const responseData = await response.json()
          setData(responseData)
          invokeSignUp(false)
        } catch (error){
          
        }
      }
      }
       
      signUp()
    }, [signUpStatus])

    console.log(document.cookie)
  
    
  return (<>
    <label htmlFor="username">Username</label>
    <input onChange={(e) => setName(e.target.value)} name="username" />
    <label htmlFor="email">Email</label>
    <input onChange={(e) => setEmail(e.target.value)} name="email"/>
    <label htmlFor="password">Password</label>
    <input onChange={(e) => setPass(e.target.value)} type="password" name="password"/>
    <button onClick={() => invokeSignUp(true)} type="submit">Submit</button>

    <button onClick={() => (navigate("/chats"))}>Chat</button>
  </>
  )
}

export default signUpForm