import Styles from "./modal.module.sass";
import IconButton from "@components/icon_button/icon-button";

import { useEffect, useId, useRef } from "preact/hooks";

import type {
	ComponentChildren,
	ComponentType,
	DialogHTMLAttributes,
	JSX,
	TargetedEvent,
	TargetedMouseEvent,
} from "preact";

export type ModalCloseReason = "backdrop" | "close-control" | "escape" | "native";

export type ModalCloseComponentProps = {
	close: () => void;
	label: string;
};

type ModalNativeProps = Omit<
	DialogHTMLAttributes<HTMLDialogElement>,
	| "aria-label"
	| "aria-labelledby"
	| "children"
	| "class"
	| "className"
	| "dangerouslySetInnerHTML"
	| "onClose"
	| "open"
	| "ref"
>;

export type ModalProps = ModalNativeProps & {
	open: boolean;
	onClose: (reason: ModalCloseReason) => void;
	title: ComponentChildren;
	subtitle?: ComponentChildren;
	children: ComponentChildren;
	closeComponent?: ComponentType<ModalCloseComponentProps>;
	closeLabel?: string;
	closeOnBackdrop?: boolean;
	class?: string;
	className?: string;
	surfaceClass?: string;
	bodyClass?: string;
};

export default function Modal({
	open,
	onClose,
	title,
	subtitle,
	children,
	closeComponent,
	closeLabel = "Zamknij modal",
	closeOnBackdrop = true,
	class: className,
	className: nativeClassName,
	surfaceClass,
	bodyClass,
	onCancel,
	onClick,
	...dialogProps
}: ModalProps): JSX.Element {
	const dialogRef = useRef<HTMLDialogElement>(null);
	const titleId = useId();
	const subtitleId = useId();

	const CloseComponent = closeComponent;
	const hasSubtitle = subtitle !== undefined && subtitle !== null && subtitle !== false;

	const rootClasses = [Styles["modal"], className, nativeClassName].filter(Boolean).join(" ");
	const surfaceClasses = [Styles["modal__surface"], surfaceClass].filter(Boolean).join(" ");
	const bodyClasses = [Styles["modal__body"], bodyClass].filter(Boolean).join(" ");

	useEffect(() => {
		const dialog = dialogRef.current;

		if (!dialog) return;

		if (open && !dialog.open) dialog.showModal();
		else if (!open && dialog.open) dialog.close()
	}, [open]);

	const handleCancel = (event: TargetedEvent<HTMLDialogElement, Event>): void => {
		onCancel?.(event);

		if (event.defaultPrevented) return;

		event.preventDefault();
		onClose("escape");
	};

	const handleClick = (event: TargetedMouseEvent<HTMLDialogElement>): void => {
		onClick?.(event);

		if (
			event.defaultPrevented
			|| !closeOnBackdrop
			|| event.target !== event.currentTarget
		) return;

		onClose("backdrop");
	};

	const handleNativeClose = (): void => {
		if (open) onClose("native");
	};

	const closeFromControl = (): void => onClose("close-control");

	return (
		<dialog
			{...dialogProps}
			ref={dialogRef}
			class={rootClasses}
			aria-labelledby={titleId}
			aria-describedby={dialogProps["aria-describedby"] ?? (hasSubtitle ? subtitleId : undefined)}
			onCancel={handleCancel}
			onClose={handleNativeClose}
			onClick={handleClick}
			data-modal=""
			data-state={open ? "open" : "closed"}
		>
			<div class={surfaceClasses} data-modal-surface="">
				<div class={Styles["modal__handle"]} aria-hidden="true" data-modal-handle="" />

				<header
					class={Styles["modal__header"]}
					data-modal-header=""
					data-has-subtitle={hasSubtitle ? "true" : undefined}
				>
					<div class={Styles["modal__heading"]} data-modal-heading="">
						<h2 id={titleId} class={Styles["modal__title"]} data-modal-title="">
							{title}
						</h2>

						{hasSubtitle && (
							<div
								id={subtitleId}
								class={Styles["modal__subtitle"]}
								data-modal-subtitle=""
							>
								{subtitle}
							</div>
						)}
					</div>

					<div class={Styles["modal__close"]} data-modal-close="">
						{CloseComponent
							? <CloseComponent close={closeFromControl} label={closeLabel} />
							: (
								<IconButton
									class={Styles["modal__default-close"]}
									icon="fa7-solid:xmark"
									aria-label={closeLabel}
									onClick={closeFromControl}
									data-modal-default-close=""
								/>
							)}
					</div>
				</header>

				<div class={bodyClasses} data-modal-body="">
					{children}
				</div>
			</div>
		</dialog>
	);
}
