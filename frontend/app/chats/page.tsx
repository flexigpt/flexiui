import type { FC } from 'react';
import { useCallback, useEffect, useRef, useState } from 'react';

import { v7 as uuidv7 } from 'uuid';

import type { ModelParams } from '@/models/aiprovidermodel';
import type { Conversation, ConversationItem, ConversationMessage } from '@/models/conversationmodel';
import { ConversationRoleEnum } from '@/models/conversationmodel';
import { type ChatOptions, DefaultChatOptions } from '@/models/settingmodel';

import { GetCompletionMessage } from '@/apis/aiprovider_helper';
import { conversationStoreAPI } from '@/apis/baseapi';
import { listAllConversations } from '@/apis/conversationstore_helper';

import ButtonScrollToBottom from '@/components/button_scroll_to_bottom';

import ChatInputField, { type ChatInputFieldHandle } from '@/chats/chat_input_field';
import ChatMessage from '@/chats/chat_message';
import ChatNavBar from '@/chats/chat_navbar';

function initConversation(title = 'New Conversation'): Conversation {
	return {
		id: uuidv7(),
		title: title.substring(0, 64),
		createdAt: new Date(),
		modifiedAt: new Date(),
		messages: [],
	};
}

function initConversationMessage(role: ConversationRoleEnum, content: string): ConversationMessage {
	const d = new Date();
	return {
		id: d.toISOString(),
		createdAt: new Date(),
		role,
		content,
	};
}

const ChatScreen: FC = () => {
	const [chat, setChat] = useState<Conversation>(initConversation());
	const [initialItems, setInitialItems] = useState<ConversationItem[]>([]);
	const [inputHeight, setInputHeight] = useState(0);
	const [streamedMessage, setStreamedMessage] = useState('');
	const [isStreaming, setIsStreaming] = useState(false);

	const chatContainerRef = useRef<HTMLDivElement>(null);
	const chatInputRef = useRef<ChatInputFieldHandle>(null);
	const conversationListRef = useRef<ConversationItem[]>([]);
	const isSubmittingRef = useRef(false);

	// Focus on mount.
	useEffect(() => {
		chatInputRef.current?.focus();
	}, []);

	// Fetch conversations / search helpers.
	const fetchConversations = useCallback(async () => {
		const conversations = await listAllConversations();
		conversationListRef.current = conversations;
		setInitialItems(conversations);
	}, []);

	useEffect(() => {
		fetchConversations();
	}, [fetchConversations]);

	const fetchSearchResults = async (query: string): Promise<ConversationItem[]> => {
		if (conversationListRef.current.length === 0) {
			await fetchConversations();
		}
		return conversationListRef.current.filter(item => item.title.toLowerCase().includes(query.toLowerCase()));
	};

	const handleNewChat = useCallback(async () => {
		if (chat.messages.length === 0) {
			chatInputRef.current?.focus();
			return;
		}
		conversationStoreAPI.saveConversation(chat);
		setChat(initConversation());
		fetchConversations();
		chatInputRef.current?.focus();
	}, [chat, fetchConversations]);

	// Scroll to bottom.
	useEffect(() => {
		if (chatContainerRef.current) {
			chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
		}
	}, [chat.messages, streamedMessage]);

	const handleSelectConversation = useCallback(async (item: ConversationItem) => {
		const selectedChat = await conversationStoreAPI.getConversation(item.id, item.title);
		if (selectedChat) setChat(selectedChat);
	}, []);

	const getConversationForExport = useCallback(async (): Promise<string> => {
		const selectedChat = await conversationStoreAPI.getConversation(chat.id, chat.title);
		return JSON.stringify(selectedChat, null, 2);
	}, [chat.id, chat.title]);

	const saveUpdatedChat = (updatedChat: Conversation) => {
		conversationStoreAPI.saveConversation(updatedChat);
		setChat(updatedChat);
	};

	const updateStreamingMessage = useCallback(async (updatedChatWithUserMessage: Conversation, options: ChatOptions) => {
		let prevMessages = updatedChatWithUserMessage.messages;
		if (options.disablePreviousMessages) {
			prevMessages = [updatedChatWithUserMessage.messages[updatedChatWithUserMessage.messages.length - 1]];
		}

		setIsStreaming(true);
		setStreamedMessage('');

		await new Promise(res => setTimeout(res, 0));

		const convoMsg = initConversationMessage(ConversationRoleEnum.assistant, '');
		const updatedChatWithConvoMessage = {
			...updatedChatWithUserMessage,
			messages: [...updatedChatWithUserMessage.messages, convoMsg],
			modifiedAt: new Date(),
		};

		const onStreamData = (data: string) => {
			setStreamedMessage(prev => {
				if (prev === '') {
					updatedChatWithConvoMessage.messages[updatedChatWithConvoMessage.messages.length - 1].content = data;
					updatedChatWithConvoMessage.modifiedAt = new Date();
					setChat(updatedChatWithConvoMessage);
					return data;
				}
				return prev + data;
			});
		};

		const inputParams: ModelParams = {
			name: options.name,
			temperature: options.temperature,
			stream: options.stream,
			maxPromptLength: options.maxPromptLength,
			maxOutputLength: options.maxOutputLength,
			reasoning: options.reasoning,
			systemPrompt: options.systemPrompt,
			timeout: options.timeout,
			additionalParameters: options.additionalParameters,
		};

		const newMsg = await GetCompletionMessage(options.provider, inputParams, convoMsg, prevMessages, onStreamData);

		if (newMsg.requestDetails) {
			if (updatedChatWithConvoMessage.messages.length > 1) {
				updatedChatWithConvoMessage.messages[updatedChatWithConvoMessage.messages.length - 2].details =
					newMsg.requestDetails;
			}
		}

		if (newMsg.responseMessage) {
			const respMessage = newMsg.responseMessage;
			updatedChatWithConvoMessage.messages.pop();
			updatedChatWithConvoMessage.messages.push(respMessage);
			updatedChatWithConvoMessage.modifiedAt = new Date();
			saveUpdatedChat(updatedChatWithConvoMessage);
			setIsStreaming(false);
		}

		isSubmittingRef.current = false;
	}, []);

	const sendMessage = async (text: string, options: ChatOptions) => {
		if (isSubmittingRef.current) return;
		isSubmittingRef.current = true;

		const trimmed = text.trim();
		if (!trimmed) {
			isSubmittingRef.current = false;
			return;
		}

		const newMsg = initConversationMessage(ConversationRoleEnum.user, trimmed);
		const updated = {
			...chat,
			messages: [...chat.messages, newMsg],
			modifiedAt: new Date(),
		};

		if (updated.messages.length === 1) {
			const t = updated.messages[0].content.substring(0, 48);
			updated.title = t.charAt(0).toUpperCase() + t.slice(1);
		}

		saveUpdatedChat(updated);
		updateStreamingMessage(updated, options);
	};

	const handleEdit = useCallback(
		async (edited: string, id: string) => {
			const idx = chat.messages.findIndex(m => m.id === id);
			if (idx === -1) return;

			const msgs = chat.messages.slice(0, idx + 1);
			msgs[idx].content = edited;

			const updated = { ...chat, messages: msgs, modifiedAt: new Date() };
			saveUpdatedChat(updated);

			let opts = { ...DefaultChatOptions };
			if (chatInputRef.current) opts = chatInputRef.current.getChatOptions();
			await updateStreamingMessage(updated, opts);
		},
		[chat, updateStreamingMessage]
	);

	const handleResend = useCallback(
		async (id: string) => {
			const idx = chat.messages.findIndex(m => m.id === id);
			if (idx === -1) return;

			const msgs = chat.messages.slice(0, idx + 1);
			const updated = { ...chat, messages: msgs, modifiedAt: new Date() };
			saveUpdatedChat(updated);

			let opts = { ...DefaultChatOptions };
			if (chatInputRef.current) opts = chatInputRef.current.getChatOptions();
			await updateStreamingMessage(updated, opts);
		},
		[chat, updateStreamingMessage]
	);

	const renderedMessages = chat.messages.map((msg, idx) => {
		const live =
			isStreaming && idx === chat.messages.length - 1 && msg.role === ConversationRoleEnum.assistant
				? streamedMessage
				: '';

		return (
			<ChatMessage
				key={msg.id}
				message={msg}
				onEdit={txt => handleEdit(txt, msg.id)}
				onResend={() => handleResend(msg.id)}
				streamedMessage={live}
			/>
		);
	});

	return (
		<div className="flex flex-col items-center w-full h-full overflow-hidden">
			{/* NAVBAR */}
			<div className="w-full flex justify-center fixed top-2 z-10">
				<div className="w-11/12 lg:w-4/5 xl:w-3/4">
					<ChatNavBar
						onNewChat={handleNewChat}
						getConversationForExport={getConversationForExport}
						initialSearchItems={initialItems}
						onSearch={fetchSearchResults}
						onSelectConversation={handleSelectConversation}
						chatTitle={chat.title}
					/>
				</div>
			</div>

			{/* MESSAGES */}
			<div className="flex flex-col items-center w-full grow overflow-hidden mt-32">
				<div
					className="w-full grow flex justify-center overflow-y-auto"
					ref={chatContainerRef}
					style={{ maxHeight: `calc(100vh - 196px - ${inputHeight}px)` }}
				>
					<div className="w-11/12 lg:w-4/5 xl:w-3/4">
						<div className="w-full flex-1 space-y-4">{renderedMessages}</div>
					</div>
				</div>

				{/* INPUT */}
				<div className="w-full flex justify-center fixed bottom-0 mb-3">
					<div className="w-11/12 lg:w-4/5 xl:w-3/4">
						<ChatInputField ref={chatInputRef} onSend={sendMessage} setInputHeight={setInputHeight} />
					</div>
				</div>
			</div>

			{/* SCROLL-TO-BOTTOM BUTTON */}
			<div className="fixed bottom-0 right-0 mb-16 mr-0 lg:mr-16 z-10">
				<ButtonScrollToBottom
					scrollContainerRef={chatContainerRef}
					size={32}
					className="btn btn-md bg-transparent border-none flex items-center shadow-none"
				/>
			</div>
		</div>
	);
};

export default ChatScreen;
