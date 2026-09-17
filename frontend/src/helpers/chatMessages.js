export function normalizeChatMessage(input) {
  const message = input?.message || input || {}
  return {
    id: Number(message.id),
    chatId: Number(message.chatId ?? message.chat_id),
    senderId: Number(message.senderId ?? message.sender_id),
    username: message.username || '',
    firstName: message.firstName || '',
    lastName: message.lastName || '',
    avatarPath: message.avatarPath || '',
    content: message.content || '',
    createdAt: message.createdAt ?? message.created_at ?? '',
    isOwn: Boolean(message.isOwn),
  }
}

export function normalizeChatMessages(messages) {
  return (messages || []).map(normalizeChatMessage)
}

export function appendUniqueMessage(messages, message) {
  if (!message.id || messages.some(item => item.id === message.id)) return messages
  return [...messages, message]
}
