import { useCallback } from "preact/hooks";
import Modal from "./modal";
import Pagination from "./pagination";
import { formatRelativeTime } from "../util/time";
import Clock from "../icons/clock.svg";

const SolvesDialog = ({
  onClose,
  challName,
  solveCount,
  solves,
  page,
  setPage,
  pageSize,
  modalBodyRef,
  ...props
}) => {
  const wrappedOnClose = useCallback(
    (e) => {
      e.preventDefault();
      onClose();
    },
    [onClose]
  );

  return (
    <Modal {...props} open={solves !== null} onClose={onClose}>
      {solves !== null && (
        <>
          {solves.length === 0 ? (
            <div>
              <img src={Clock} style={{ filter: "invert(1)" }} />
              <h5>{challName} has no solves.</h5>
            </div>
          ) : (
            <>
              <div class="modal-header">
                <div class="modal-title">Solves for {challName}</div>
              </div>
              <div class={`modal-body`} ref={modalBodyRef}>
                <div>
                  <div>#</div>
                  <div>Team</div>
                  <div>Solve time</div>
                  {solves.map((solve, i) => (
                    <>
                      #{(page - 1) * pageSize + i + 1}Team
                      <a href={`/profile/${solve.userId}`}>{solve.userName}</a>
                      Solve time
                      <div>{formatRelativeTime(solve.createdAt)}</div>
                    </>
                  ))}
                </div>
                <Pagination
                  {...{ totalItems: solveCount, pageSize, page, setPage }}
                  numVisiblePages={9}
                />
              </div>
            </>
          )}
          <div class="modal-footer">
            <div class="btn-container u-inline-block">
              <button onClick={wrappedOnClose}>Close</button>
            </div>
          </div>
        </>
      )}
    </Modal>
  );
};

export default SolvesDialog;
