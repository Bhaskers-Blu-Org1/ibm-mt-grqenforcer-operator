FROM scratch

#IBMDEV Switch to ibm-mt-grqenforcer-operator binary
COPY ibm-mt-grqenforcer-operator .
CMD ["./ibm-mt-grqenforcer-operator"]
