import type { CSSProperties } from 'react';

import { useModalDialogController } from '@/hooks/use_dialog_controller';

import { ModalBackdrop } from '@/components/modal/modal_backdrop';
import { ModalDialog } from '@/components/modal/modal_dialog';

interface MermaidZoomModalProps {
	isOpen: boolean;
	onClose: () => void;
	imageSrc: string;
	surfaceStyle?: CSSProperties;
}

function MermaidZoomModalContent({ imageSrc, surfaceStyle }: Omit<MermaidZoomModalProps, 'isOpen' | 'onClose'>) {
	const { requestClose } = useModalDialogController();

	return (
		<>
			<button
				type="button"
				className="modal-box app-bg-mermaid flex h-[90vh] max-w-[90vw] cursor-zoom-out items-center justify-center border-0"
				style={surfaceStyle}
				aria-label="Close enlarged Mermaid diagram"
				onClick={() => {
					requestClose();
				}}
			>
				<img
					src={imageSrc}
					alt=""
					aria-hidden="true"
					decoding="async"
					draggable={false}
					className="pointer-events-none block size-auto max-h-[80vh] max-w-[86vw] object-contain"
				/>
			</button>

			<ModalBackdrop enabled={true} />
		</>
	);
}

export function MermaidZoomModal(props: MermaidZoomModalProps) {
	if (!props.isOpen) {
		return null;
	}

	return (
		<ModalDialog isOpen={props.isOpen} onClose={props.onClose} aria-label="Enlarged Mermaid diagram">
			<MermaidZoomModalContent imageSrc={props.imageSrc} surfaceStyle={props.surfaceStyle} />
		</ModalDialog>
	);
}
