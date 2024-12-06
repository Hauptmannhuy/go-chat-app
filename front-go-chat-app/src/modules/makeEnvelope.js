export const useEnvelope = () => {
  const makeEnvelope = (type, data) => {
    const actions = {
      "NEW_MESSAGE": {
        type: "NEW_MESSAGE",
        chat_name: `${data[0]}`,
        body: `${data[1]}`,
      },
      "NEW_CHAT": {
        type: "NEW_CHAT",
        chat_name: `${data[0]}`,
      },
      "NEW_PRIVATE_CHAT": {
        type: "NEW_PRIVATE_CHAT",
        receiver_id: `${data[0]}`,
        message: `${data[1]}`,
      },
      "JOIN_CHAT": {
        type: "JOIN_CHAT",
        chat_id: `${data[0]}`,
      },
      "SEARCH_QUERY": {
        type: "SEARCH_QUERY",
        input: `${data[0]}`,
      }
    }

    const json = actions[type];
    if (json) {
      return json
    } else {
      console.error(`Unknown action type: ${type}`);
    }
  };
  return {makeEnvelope}
}