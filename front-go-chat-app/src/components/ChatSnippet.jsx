export function ChatSnippet({messages, name, currentUsername}) {
  const nameFormatter = (name) => (name.split('_').filter((el) => (el != currentUsername)).join(' '))
  const dialogue = messages[name]
  const lastMessage = dialogue[dialogue.length-1]
  if (lastMessage) {
    return (
      <div className="snippet-container">
        <div className="chatname-snippet"> {nameFormatter(name)}</div>
        <div className="chatmessage-snippet"> {lastMessage.body}</div>
      </div>
      )
    }  else {
      return (
        <div className="snippet-container">
          <div className="chatname-snippet">{nameFormatter(name)}</div>
        </div>
      )
    }
}