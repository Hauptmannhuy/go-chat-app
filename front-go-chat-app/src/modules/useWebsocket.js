import { useRef } from "react";
import { useNavigate } from "react-router-dom";

export function useWebsocket(url, onMessage) {
  
  const socket = useRef(null)
  const navigate = useNavigate()

  
  async function connectWS() {
      checkAuth().then(() => {
      const websocket =  new WebSocket(url)
      websocket.addEventListener("open", () => {
      socket.current = websocket
      console.log("connected to WS")
      })
      websocket.addEventListener("message", (ev) => {
        onMessage(ev)
      })
      websocket.addEventListener("error", (ev) => {
        console.log(ev)
        // navigate("/sign_up")
        throw new Error("Error connecting to WebSocket");
      })
    })
    .catch((reason) => {
      // navigate("/sign_up")
      console.log(reason)
      return new Error(reason)
    })
  }


  const checkAuth = async() => {
    try {
     const response = await fetch("api/checkauth", {
        mode: 'cors',
        credentials: 'include'
      })
      if(response.status == 200 ){
       return "ok"
      } else {
       throw new Error("not authorized");
      }
    } catch (error) {
      console.error("Error checking auth:", error)
      throw new Error(error);
    }
 }
     
 const writeToSocket = (message) => {
  if (socket.current){
   socket.current.send(JSON.stringify(message))
 } else {
   setTimeout(() => {
    writeToSocket(message)
   }, 1000);
 }
}


 
  return { connectWS, writeToSocket, socket}
}