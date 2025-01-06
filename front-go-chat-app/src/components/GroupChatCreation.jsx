import { useContext, useState } from "react"
import { useEnvelope } from "../modules/useEnvelope"
import { GlobalContext } from "../contexts/GlobalContext"


export default function CreateGroupChat(){
  const [creationChatInvoked, setInvokeStatus] = useState(false)
  const [creationInput, setCreationInput] = useState("")
  const {writeToSocket} = useContext(GlobalContext)
  const {makeEnvelope} = useEnvelope()

  const createGroupChat = (groupName) => {
    const message = makeEnvelope("NEW_GROUP_CHAT", [groupName] )
    writeToSocket(message)
    setInvokeStatus(false)
    setCreationInput("")
  }

  return (
    <>
      <button onClick={() => setInvokeStatus(true)}>Create Group Chat</button>
      { (creationChatInvoked) ? (
        <div>
          <input type="text" onChange={(e) => setCreationInput(e.target.value)} /> 
          <button onClick={() => createGroupChat(creationInput)}>Submit</button>
        </div>
      )
      :
      (
        null
      )}
    </>
  )
}