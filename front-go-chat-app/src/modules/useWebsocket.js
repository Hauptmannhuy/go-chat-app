import { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";


export function useWebsocket(url, onMessage) {
  
  const socket = useRef(null)
  const [statusCode, setStatusCode] = useState(null)
  const navigate = useNavigate()
  
  useEffect(() => {

      if (statusCode != null){
        if (statusCode == 200) {

          const websocket =  new WebSocket(url)
          
          websocket.addEventListener("open", () => {
            socket.current = websocket
          })
          
          websocket.addEventListener("message", onMessage)
          
          return () => websocket.close()
        } else {
          navigate("/sign_up")
        }
    } 
  }, [statusCode])

  
   const sendMessage = (message) => {
     if (socket.current){
      console.log(socket)
      socket.current.send(JSON.stringify(message))
    } else {
      setTimeout(() => {
        sendMessage(message)
      }, 1000);
    }
  }


  useEffect(() => {
    const checkAuth = async() => {
     try {
      const response = await fetch("api/checkauth", {
         mode: 'cors',
         credentials: 'include'
       })

       setStatusCode(response.status)
     } catch (error) {
       console.log(error)
     }
     }
     checkAuth()
   }, [])
 
  return { sendMessage, socket}
}