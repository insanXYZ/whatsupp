import { toast } from "sonner";

export function ToastSuccess(message: string) {
  toast.success(message);
}

export function ToastInfo(message: string) {
  toast.info(message);
}

export function ToastError(message: string, description?: string) {
  toast.error(message, {
    description: description,
  });
}
