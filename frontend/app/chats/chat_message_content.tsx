import type { ReactNode } from 'react';
import { memo, useEffect, useMemo, useState } from 'react';

import 'katex/dist/katex.min.css';
import Markdown from 'react-markdown';
import rehypeKatex from 'rehype-katex';
import remarkGemoji from 'remark-gemoji';
import remarkGfm from 'remark-gfm';
import remarkMath from 'remark-math';
import supersub from 'remark-supersub';

import { backendAPI } from '@/apis/baseapi';

import CodeBlock from '@/chats/chat_message_content_codeblock';

const remarkPlugins = [remarkGemoji, supersub, remarkMath, remarkGfm];
const rehypePlugins = [rehypeKatex];
const remarkPluginsStreaming = [remarkGemoji, supersub, remarkGfm];

function useDebounced<T>(value: T, delay: number): T {
	const [debounced, setDebounced] = useState(value);
	useEffect(() => {
		const id = setTimeout(() => {
			setDebounced(value);
		}, delay);
		return () => {
			clearTimeout(id);
		};
	}, [value, delay]);
	return debounced;
}

// LaTeX processing function
const containsLatexRegex = /\\\(.*?\\\)|\\\[.*?\\\]|\$.*?\$|\\begin\{equation\}.*?\\end\{equation\}/;
const inlineLatex = new RegExp(/\\\((.+?)\\\)/, 'g');
const blockLatex = new RegExp(/\\\[(.*?[^\\])\\\]/, 'gs');

const processLaTeX = (content: string) => {
	let processedContent = content.replace(/(\$)(?=\s?\d)/g, '\\$');

	if (!containsLatexRegex.test(processedContent)) {
		return processedContent;
	}

	processedContent = processedContent
		.replace(inlineLatex, (match: string, equation: string) => `$${equation}$`)
		.replace(blockLatex, (match: string, equation: string) => `$$${equation}$$`);

	return processedContent;
};

interface CodeComponentProps {
	inline?: boolean;
	className?: string;
	children?: ReactNode;
}

interface PComponentProps {
	className?: string;
	children?: ReactNode;
}

interface ChatMessageContentProps {
	content: string; // final text (when stream finished)
	streamedText?: string; // partial text while streaming
	isStreaming?: boolean;
	align: string;
	renderAsMarkdown?: boolean;
}

const ChatMessageContentBase = ({
	content,
	streamedText = '',
	isStreaming = false,
	align,
	renderAsMarkdown = true,
}: ChatMessageContentProps) => {
	const liveText = isStreaming ? streamedText : content;
	const textToRender = useDebounced(liveText, 250); // parse max ~4×/sec

	if (!renderAsMarkdown) {
		// Memoize the plain text content to prevent unnecessary re-renders
		const plainTextContent = useMemo(() => {
			return textToRender.split('\n').map((line, index) => (
				<p
					key={index}
					className={`${align} break-words`}
					style={{ whiteSpace: 'pre-wrap', lineHeight: '1.5', fontSize: '14px' }}
				>
					{line || '\u00A0' /* Use non-breaking space for empty lines */}
				</p>
			));
		}, [textToRender, align]);
		return <div className="bg-base-100 px-4 py-2">{plainTextContent}</div>;
	}

	const processedContent = useMemo(() => {
		if (isStreaming) return textToRender; // skip LaTeX processing
		return /[$\\]/.test(textToRender) ? processLaTeX(textToRender) : textToRender;
	}, [textToRender, isStreaming]);

	const components = useMemo(
		() => ({
			h1: ({ children }: PComponentProps) => <h1 className="text-xl font-bold my-2">{children}</h1>,
			h2: ({ children }: PComponentProps) => <h2 className="text-lg font-bold my-2">{children}</h2>,
			h3: ({ children }: PComponentProps) => <h3 className="text-base font-bold my-2">{children}</h3>,
			p: ({ className, children }: PComponentProps) => {
				const newClassName = `${className || ''} my-2 ${align} break-words`;
				return (
					<p className={newClassName} style={{ lineHeight: '1.5', fontSize: '14px' }}>
						{children}
					</p>
				);
			},
			code: ({ inline, className, children, ...props }: CodeComponentProps) => {
				if (inline || !className) {
					const newClassName = `bg-base-200 inline text-wrap whitespace-pre-wrap break-words ${className}`;
					return (
						<code className={newClassName} {...props}>
							{children}
						</code>
					);
				}
				const match = /lang-(\w+)/.exec(className || '') || /language-(\w+)/.exec(className || '');
				const language = match && match[1] ? match[1] : 'text';

				return (
					<CodeBlock
						language={language}
						// eslint-disable-next-line @typescript-eslint/no-base-to-string
						value={String(children).replace(/\n$/, '')}
						streamedMessage={streamedText}
						{...props}
					/>
				);
			},
			ul: ({ children }: PComponentProps) => (
				<span>
					<ul className="list-disc py-1">{children}</ul>
				</span>
			),
			ol: ({ children }: PComponentProps) => (
				<span>
					<ol className="list-decimal py-1">{children}</ol>
				</span>
			),
			li: ({ children }: PComponentProps) => (
				<span>
					<li className="ml-4 py-1">{children}</li>
				</span>
			),
			table: ({ children }: PComponentProps) => <table className="table-auto w-full">{children}</table>,
			thead: ({ children }: PComponentProps) => <thead className="bg-base-300">{children}</thead>,
			tbody: ({ children }: PComponentProps) => <tbody>{children}</tbody>,
			tr: ({ children }: PComponentProps) => <tr className="border-t">{children}</tr>,
			th: ({ children }: PComponentProps) => <th className="px-4 py-2 text-left">{children}</th>,
			td: ({ children }: PComponentProps) => <td className="px-4 py-2">{children}</td>,
			a: ({ href, children }: { href?: string; children?: ReactNode }) => (
				<a
					href={href}
					target={href?.startsWith('http') ? '_blank' : undefined}
					rel={href?.startsWith('http') ? 'noopener noreferrer' : undefined}
					className="underline text-blue-600 hover:text-blue-800 cursor-pointer"
					onClick={e => {
						e.preventDefault();
						if (href) {
							backendAPI.openurl(href);
						}
					}}
				>
					{children}
				</a>
			),
			blockquote: ({ children }: PComponentProps) => (
				<blockquote className="border-l-4 border-neutral/20 pl-4 italic">{children}</blockquote>
			),
		}),
		[align]
	);

	return (
		<div className="bg-base-100 px-4 py-2">
			<Markdown
				remarkPlugins={isStreaming ? remarkPluginsStreaming : remarkPlugins}
				rehypePlugins={isStreaming ? [] : rehypePlugins}
				components={components}
			>
				{processedContent}
			</Markdown>
		</div>
	);
};

function contentAreEqual(prev: ChatMessageContentProps, next: ChatMessageContentProps) {
	return (
		prev.content === next.content &&
		prev.streamedText === next.streamedText &&
		prev.isStreaming === next.isStreaming &&
		prev.align === next.align
	);
}

export default memo(ChatMessageContentBase, contentAreEqual);
