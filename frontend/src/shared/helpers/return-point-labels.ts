import type {
	OpeningHours,
	ReturnPointType,
	ReturnPointStatus,
	ReturnPointMaterial,
	CommunityReportStatus,
	ReturnPointRefundMethod,
} from "@shared/api/responses";

export type PointStatus = ReturnPointStatus | CommunityReportStatus;

export const POINT_STATUS_LABELS: Record<PointStatus, string> = {
	dostepny: "Dostępny",
	nieczynny: "Nieczynny",
	brak_informacji: "Brak informacji",
};

export const RETURN_POINT_TYPE_LABELS: Record<ReturnPointType, string> = {
	butelkomat: "Butelkomat",
	sklep: "Sklep",
	punkt_skupu: "Punkt skupu",
	inny: "Inny punkt",
};

export const RETURN_POINT_MATERIAL_LABELS: Record<ReturnPointMaterial, string> = {
	pet: "Butelki PET",
	can: "Puszki metalowe",
	glass: "Butelki szklane",
};

export const RETURN_POINT_MATERIAL_SHORT_LABELS: Record<ReturnPointMaterial, string> = {
	pet: "PET",
	can: "Puszki",
	glass: "Szklane",
};

export const RETURN_POINT_REFUND_METHOD_LABELS: Record<ReturnPointRefundMethod, string> = {
	gotowka: "Gotówka",
	bon: "Bon",
	karta: "Karta",
};

export const OPENING_HOURS_LABELS: Record<keyof OpeningHours, string> = {
	"mon-fri": "Pon–pt",
	"mon-sat": "Pon–sob",
	"mon-sun": "Pon–niedz",
	sat: "Sobota",
	sun: "Niedziela",
};
