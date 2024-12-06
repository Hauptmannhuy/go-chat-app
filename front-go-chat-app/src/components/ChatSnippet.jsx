export function ChatSnippet({messages, name, currentUsername}) {
  const nameFormatter = (name) => (name.split('_').filter((el) => (el != currentUsername)).join(' '))
  const dialogue = messages[name]
  // const lastMessage = dialogue[dialogue.length-1]
  if (dialogue && dialogue.length > 0) {
    return (
      <div className="snippet-container">
        <div className="chatname-snippet"> {nameFormatter(name)}</div>
        <div className="chatmessage-snippet"> {dialogue[dialogue.length-1].body}</div>
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