import Styles from "./navigation.module.sass";
import IconifyIcon from "@components/iconify_icon/iconify-icon";

import { useEffect, useId, useRef, useState } from "preact/hooks";

export type NavigationSheetLink = {
	icon: string;
	href: string;
	label: string;
	active: boolean;
};

type NavigationSheetProps = {
	links: NavigationSheetLink[];
	active?: boolean;
};

export default function NavigationSheet({ links, active = false }: NavigationSheetProps) {
	const [isOpen, setIsOpen] = useState(false);

	const dialogRef = useRef<HTMLDialogElement>(null);

	const titleId = useId();
	const dialogId = useId();

	const openSheet = (): void => {
		const dialog = dialogRef.current;

		if (!dialog || dialog.open) return;

		dialog.showModal();

		setIsOpen(true);
	};

	const closeSheet = (): void => {
		const dialog = dialogRef.current;

		if (dialog?.open) dialog.close();
	};

	useEffect(() => {
		const desktopMedia = window.matchMedia("(min-width: 1024px)");

		const closeOnDesktop = (event: MediaQueryListEvent): void => {
			if (event.matches && dialogRef.current?.open) dialogRef.current.close();
		};

		desktopMedia.addEventListener("change", closeOnDesktop);

		return () => desktopMedia.removeEventListener("change", closeOnDesktop);
	}, []);

	return (
		<div
			class={Styles["navigation-sheet"]}
			data-active={active ? "true" : undefined}
			data-state={isOpen ? "open" : "closed"}
		>
			<button
				type="button"
				class={`${Styles["navigation-sheet__trigger"]} ${
					active ? Styles["navigation-sheet__trigger--active"] : ""
				}`}
				aria-controls={dialogId}
				aria-expanded={isOpen}
				aria-haspopup="dialog"
				onClick={openSheet}
			>
				<span class={Styles["navigation-sheet__indicator"]} aria-hidden="true">
					<IconifyIcon
						class={Styles["navigation-sheet__trigger-icon"]}
						icon="fa7-solid:ellipsis"
					/>
				</span>

				<span class={Styles["navigation-sheet__trigger-label"]}>Więcej</span>
			</button>

			<dialog
				ref={dialogRef}
				id={dialogId}
				class={Styles["navigation-sheet__dialog"]}
				aria-labelledby={titleId}
				onClose={() => setIsOpen(false)}
				onClick={(event) => {
					if (event.target === event.currentTarget) closeSheet();
				}}
			>
				<div class={Styles["navigation-sheet__handle"]} aria-hidden="true" />

				<header class={Styles["navigation-sheet__header"]}>
					<h2 id={titleId} class={Styles["navigation-sheet__title"]}>Więcej</h2>

					<button
						type="button"
						class={Styles["navigation-sheet__close"]}
						aria-label="Zamknij menu"
						onClick={closeSheet}
					>
						<IconifyIcon icon="fa7-solid:xmark" />
					</button>
				</header>

				<ul class={Styles["navigation-sheet__list"]}>
					{links.map(({ href, label, icon, active: linkActive }) => (
						<li key={`${href}:${label}`} class={Styles["navigation-sheet__item"]}>
							<a
								href={href}
								class={`${Styles["navigation-sheet__link"]} ${
									linkActive ? Styles["navigation-sheet__link--active"] : ""
								}`}
								aria-current={linkActive ? "page" : undefined}
								data-active={linkActive ? "true" : undefined}
								onClick={closeSheet}
							>
								<span class={Styles["navigation-sheet__link-icon"]} aria-hidden="true">
									<IconifyIcon icon={icon} />
								</span>

								<span>{label}</span>
							</a>
						</li>
					))}
				</ul>
			</dialog>
		</div>
	);
}
