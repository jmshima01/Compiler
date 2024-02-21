CXX = g++
LD = g++
CXXFLAGS = -g3 -std=c++17
LDFLAGS = 

#
# You should be able to add object files here without changing anything else
#
TARGET = luther
OBJ_FILES = ${TARGET}.o
INC_FILES = ${TARGET}.hh

#
# Any libraries we might need.
#
# LIBRARYS = -lbsd -lboost_regex -lboost_log_setup -lboost_log -lboost_thread -lboost_system -lpthread -lboost_filesystem

${TARGET}: ${OBJ_FILES}
	${LD} ${LDFLAGS} ${OBJ_FILES} -o $@ ${LIBRARYS}

%.o : %.cc ${INC_FILES}
	${CXX} -c ${CXXFLAGS} -o $@ $<

clean:
	rm -f core ${TARGET} ${OBJ_FILES}

check:
	valgrind --leak-check=full ./${TARGET}
