export default function SignUpInput({actionSubmit, nameHandler, emailHandler, passHandler }) {

  return ( <>
    <label htmlFor="username">Username</label>
    <input onChange={(e) => nameHandler(e.target.value)} name="username" />
    <label htmlFor="email">Email</label>
    <input onChange={(e) => emailHandler(e.target.value)} name="email"/>
    <label htmlFor="password">Password</label>
    <input onChange={(e) => passHandler(e.target.value)} type="password" name="password"/>
    <button onClick={() => actionSubmit(true)} type="submit">Submit</button>
  
  
  </>)
}